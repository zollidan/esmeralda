from django.contrib import admin

from django.contrib import admin
from .models import ParsingTask
from .tasks import run_parsing_task

@admin.register(ParsingTask)
class ParsingTaskAdmin(admin.ModelAdmin):
    list_display = ('status', 'created_at', 'result_file')
    list_filter = ('status',)
    actions = ['start_parsing']

    def start_parsing(self, request, queryset):
        for task in queryset:
            run_parsing_task.delay(task.id)
        self.message_user(request, "Задачи отправлены воркеру")
    start_parsing.short_description = "Запустить парсинг выбранных"
