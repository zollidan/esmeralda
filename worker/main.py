from enum import Enum
import json
import time
from typing import Optional
from pydantic import BaseModel
import redis

from worker.parsers import test_parser, test_parser2

PARSERS_KEY = "available_parsers"
QUEUE = "tasks_queue"

class TaskMessage(BaseModel):
    parser_name: str
    task_id: str

class TaskStatus(str, Enum):
    SUCCESS = "success"
    ERROR = "error"

class TaskResult(BaseModel):
    status: TaskStatus
    result: Optional[str] = None
    error: Optional[str] = None

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

    print("Worker started, waiting for tasks...")

    while True:
        task = r.blpop(QUEUE)
        task_data = json.loads(task[1])
        parser_name = task_data["parser_name"]
        print("Got task:", parser_name)

        try:
            result = parsers[parser_name](task_data["task_id"])
            response = TaskResult(status=TaskStatus.SUCCESS, result=result)
            r.set(f"result:{task_data['task_id']}", response.model_dump_json())
        except Exception as e:
            response = TaskResult(status=TaskStatus.ERROR, error=str(e))
            r.set(f"result:{task_data['task_id']}", response.model_dump_json())

if __name__ == "__main__":
    main()