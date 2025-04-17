import { Component, OnInit } from '@angular/core';
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
  constructor(private messageService: MessageService,public productService: ProductService) {}

  ngOnInit() {

    this.messageService.currentMessage.subscribe(
      (message: string) => this.message = message 
    );
  }
}
