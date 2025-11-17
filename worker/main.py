from enum import Enum
import json
import time
import redis

PARSERS_KEY = "available_parsers"
STREAM_NAME = "parser_tasks"
GROUP_NAME = "workers"
CONSUMER_NAME = "worker-1"


def test_parser(url):
    return f"Parsed data from {url} using test_parser"


def test_parser2(url):
    return f"Parsed data from {url} using test_parser2"


class TaskStatus(str, Enum):
    SUCCESS = "success"
    ERROR = "error"


parsers = {
    "test_parser": test_parser,
    "test_parser2": test_parser2
}


def main():
    r = redis.Redis(host="localhost", port=6379, decode_responses=True)

    r.delete(PARSERS_KEY)
    parser_list = list(parsers.keys())
    r.sadd(PARSERS_KEY, *parser_list)

    print(f"Loaded parsers: {parser_list}")

    # Создать группу потребителей, если она не существует
    try:
        r.xgroup_create(STREAM_NAME, GROUP_NAME, id="0", mkstream=True)
    except redis.exceptions.ResponseError as e:
        if "BUSYGROUP" not in str(e):
            raise

    print("Worker started, waiting for tasks...")

    while True:
        tasks = r.xreadgroup(
            GROUP_NAME,
            CONSUMER_NAME,
            {STREAM_NAME: '>'},
            count=1,
            block=5000
        )

        if tasks:
            for stream, messages in tasks:
                for message_id, message_data in messages:
                    parser_name = message_data["parser"]
                    file_name = message_data["fileName"]
                    task_id = message_data["task"]
                    print("Got task:", parser_name, file_name, task_id)

                    try:
                        result = parsers[parser_name](file_name)
                        response = {
                            "task_id": task_id,
                            "status": TaskStatus.SUCCESS.value,
                            "result": result}
                        r.set(
                            f"result:{task_id}", json.dumps(response))
                    except Exception as e:
                        response = {
                            "task_id": task_id,
                            "status": TaskStatus.ERROR.value,
                            "error": str(e)}
                        r.set(
                            f"result:{task_id}", json.dumps(response))

                    r.xack(STREAM_NAME, GROUP_NAME, message_id)

        time.sleep(1)


if __name__ == "__main__":
    main()
