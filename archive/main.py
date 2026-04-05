from datetime import datetime
from typing import List, Optional

from contextlib import asynccontextmanager
from sqlalchemy import Column, Integer, String, DateTime, select
from sqlalchemy.ext.asyncio import AsyncSession, create_async_engine, async_sessionmaker
from sqlalchemy.orm import declarative_base

from fastapi import FastAPI, Depends, HTTPException, status
from pydantic import BaseModel, ConfigDict


DATABASE_URL = 'sqlite+aiosqlite:///./archive.db'

engine = create_async_engine(DATABASE_URL, connect_args={'check_same_thread': False})
AsyncSessionLocal = async_sessionmaker(bind=engine, expire_on_commit=False)
Base = declarative_base()


class ArchiveItemModel(Base):
    __tablename__ = 'archive_items'

    id = Column(Integer, primary_key=True, index=True)
    title = Column(String, nullable=False)
    created_at = Column(DateTime, default=datetime.now(), nullable=False)


class ArchiveItemCreate(BaseModel):
    title: str


class ArchiveItemUpdate(BaseModel):
    title: Optional[str] = None


class ArchiveItemRead(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    id: int
    title: str
    created_at: datetime


async def get_db() -> AsyncSession:
    async with AsyncSessionLocal() as session:
        yield session


@asynccontextmanager
async def lifespan(app: FastAPI):
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
    yield


app = FastAPI(lifespan=lifespan)


@app.post('/items', response_model=ArchiveItemRead, status_code=status.HTTP_201_CREATED)
async def create_item(
    item: ArchiveItemCreate, db: AsyncSession = Depends(get_db)
) -> ArchiveItemRead:
    new_item = ArchiveItemModel(title=item.title)
    db.add(new_item)
    await db.commit()
    await db.refresh(new_item)
    return new_item


@app.get('/items', response_model=List[ArchiveItemRead])
async def list_items(db: AsyncSession = Depends(get_db)) -> List[ArchiveItemRead]:
    result = await db.execute(select(ArchiveItemModel))
    return result.scalars().all()


@app.get('/items/{item_id}', response_model=ArchiveItemRead)
async def get_item(item_id: int, db: AsyncSession = Depends(get_db)) -> ArchiveItemRead:
    result = await db.execute(
        select(ArchiveItemModel).where(ArchiveItemModel.id == item_id)
    )
    item = result.scalar_one_or_none()
    if not item:
        raise HTTPException(status_code=404, detail='Item not found')
    return item


@app.put('/items/{item_id}', response_model=ArchiveItemRead)
async def update_item(
    item_id: int, item_data: ArchiveItemUpdate, db: AsyncSession = Depends(get_db)
) -> ArchiveItemRead:
    result = await db.execute(
        select(ArchiveItemModel).where(ArchiveItemModel.id == item_id)
    )
    item = result.scalar_one_or_none()
    if not item:
        raise HTTPException(status_code=404, detail='Item not found')

    if item_data.title is not None:
        item.title = item_data.title

    await db.commit()
    await db.refresh(item)
    return item


@app.delete('/items/{item_id}', status_code=status.HTTP_204_NO_CONTENT)
async def delete_item(item_id: int, db: AsyncSession = Depends(get_db)) -> None:
    result = await db.execute(
        select(ArchiveItemModel).where(ArchiveItemModel.id == item_id)
    )
    item = result.scalar_one_or_none()
    if not item:
        raise HTTPException(status_code=404, detail='Item not found')

    await db.delete(item)
    await db.commit()
