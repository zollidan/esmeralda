from art import tprint
import pandas as pd
import os

def soccerway(time_date: str):
    
    file_name = 'soccerway.xlsx'
    
    data = [[1,2,3], [4,5,6], [time_date]]
    
    # df = pd.DataFrame(data=data)
    
    # df.to_excel(file_name, index=False)

    # os.remove(file_name)

    return f'soccerway worked with date {time_date}' 