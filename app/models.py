from django.db import models

class ParsingTask(models.Model):
    STATUS_CHOICES = [
        ('PENDING', 'Ожидает'),
        ('PROCESSING', 'В работе'),
        ('DONE', 'Готово'),
        ('ERROR', 'Ошибка'),
    ]

    status = models.CharField(max_length=20, choices=STATUS_CHOICES, default='PENDING')
    result_file = models.FileField(upload_to='results/', null=True, blank=True)
    error_log = models.TextField(blank=True)
    created_at = models.DateTimeField(auto_now_add=True)

    def __str__(self):
        return f"{self.created_at} - {self.status}"