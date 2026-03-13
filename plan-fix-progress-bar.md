1.  Redis Streams гарантирует порядок — сообщения читаются в порядке добавления, это не проблема.
2.  Главная проблема — broadcastProgress с select { case ch <- p: default: } дропает сообщения если канал заполнен (буфер 64). При 20 горутинах которые все завершаются почти одновременно — пачка сообщений летит разом и часть теряется.
3.  WebSocket хендлер — покажи его, но скорее всего там тоже есть проблема с тем что он не читает историю из стрима при подключении, т.е. если клиент подключился позже — он пропустил часть прогресса.
    Основной фикс — не дропать сообщения и читать из Redis Stream напрямую в WS хендлере:

    ```go

    // вместо pub/sub через каналы — каждый WS читает стрим сам
    func (h *Handler) WsProgress(c *gin.Context) {
    taskID := c.Query("task_id")

        conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
        if err != nil {
            return
        }
        defer conn.Close()

        ctx, cancel := context.WithCancel(c.Request.Context())
        defer cancel()

        lastID := "0" // читаем с начала — не пропустим ничего что уже было

        for {
            // читаем напрямую из Redis Stream
            streams, err := h.rdb.XRead(ctx, &redis.XReadArgs{
                Streams: []string{queue.StreamProgress, lastID},
                Count:   100,
                Block:   2 * time.Second,
            }).Result()
            if err != nil {
                if errors.Is(err, context.Canceled) || errors.Is(err, redis.Nil) {
                    return
                }
                continue
            }

            for _, stream := range streams {
                for _, msg := range stream.Messages {
                    lastID = msg.ID

                    payload, ok := msg.Values["payload"].(string)
                    if !ok {
                        continue
                    }

                    progress, err := queue.Unmarshal[queue.TaskProgress]([]byte(payload))
                    if err != nil {
                        continue
                    }

                    // фильтруем только нужную таску
                    if progress.TaskID != taskID {
                        continue
                    }

                    if err := conn.WriteJSON(progress); err != nil {
                        return
                    }

                    if slices.Contains([]string{"done", "error", "failed", "cancelled"}, string(progress.Status)) {
                        return
                    }
                }
            }
        }

    }

    ```

    Это решает сразу все проблемы: не дропает сообщения, читает историю при позднем подключении, и убирает всю логику subscribeProgress/unsubscribeProgress/broadcastProgress из хендлера — она больше не нужна.
    Но фильтрация if progress.TaskID != taskID внутри цикла неэффективна если тасок много. Покажи текущий WS хендлер — возможно стоит сделать отдельный стрим per-task или использовать Redis Pub/Sub вместо Streams для прогресса.

    ```

    ```
