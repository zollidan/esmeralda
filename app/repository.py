def create_file_record(file, db):
    db.add(file)
    db.commit()
    db.refresh(file)