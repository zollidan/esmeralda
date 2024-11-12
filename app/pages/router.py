from fastapi import APIRouter, Depends, Request
from fastapi.templating import Jinja2Templates
from starlette.responses import FileResponse

from app.api.router import get_files 

pages_router = APIRouter(tags=['Фронтенд'])
templates = Jinja2Templates(directory='app/templates')

@pages_router.get('/')
async def get_home_page(request: Request):
    return templates.TemplateResponse(name='home.html', context={'request': request})

@pages_router.get('/parsers')
async def get_parsers_page(request: Request):
    return templates.TemplateResponse(name='parsers.html', context={'request': request})


@pages_router.get('/files')
async def get_files_page(request: Request, files=Depends(get_files)):
    
    last_file = files[-0]
    
    print(files)
    
    
    return templates.TemplateResponse(name='files.html', context={'request': request, "files":files, "last_file": last_file})