from celery import Celery

app = Celery('parser', broker='redis://redis:6379/0')  # hostname â docker-compose

@app.task
def run_soccerway_1():
    from parsers.soccerway_1.parser import main
    main()
