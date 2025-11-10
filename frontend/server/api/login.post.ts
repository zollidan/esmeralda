import z from "zod";

const bodySchema = z.object({
  username: z.string(),
  password: z.string(),
});

export default defineEventHandler(async (event) => {
  const runtimeConfig = useRuntimeConfig();

  // Валидируем тело запроса
  const { username, password } = await readValidatedBody(
    event,
    bodySchema.parse
  );

  try {
    const response = await fetch(
      `${runtimeConfig.apiBaseURL}/api/v1/auth/login`,
      {
        method: "POST",
        body: JSON.stringify({ username, password }),
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    // 💡 Обработка разных кодов ответа
    if (response.status === 401) {
      return sendError(
        event,
        createError({
          statusCode: 401,
          statusMessage: "Unauthorized: Invalid username or password",
        })
      );
    }

    if (response.ok) {
      const data = await response.json();
      await setUserSession(event, {
        user: {
          name: username,
          token: data.token,
        },
      });
    }

    // Если код ответа не 200 и не 401
    return sendError(
      event,
      createError({
        statusCode: response.status,
        statusMessage: `Login failed with status ${response.status}`,
      })
    );
  } catch (error) {
    // Обработка сетевых ошибок или неожиданных исключений
    console.error("Login error:", error);
    return sendError(
      event,
      createError({
        statusCode: 500,
        statusMessage: "Internal server error during login request",
      })
    );
  }
});
