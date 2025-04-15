import { Component, OnInit } from '@angular/core';
import { ApiService } from './services/api.service'; // Asegúrate que la ruta es correcta
import { MessageService } from './message.service';
import { ProductService } from './product.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  message: string = '';
  title: string = 'Gestión de Productos';
  tasks: any[] = []; // Nueva propiedad para almacenar tareas

  constructor(
    private messageService: MessageService,
    public productService: ProductService,
    private apiService: ApiService // Inyecta el nuevo servicio
  ) {}

  ngOnInit() {
    this.messageService.currentMessage.subscribe(
      (message: string) => this.message = message
    );
    
    // Cargar tareas al iniciar
    this.loadTasks();
  }

  loadTasks() {
    this.apiService.getTasks().subscribe({
      next: (tasks) => this.tasks = tasks,
      error: (err) => console.error('Error cargando tareas:', err)
    });
  }
}