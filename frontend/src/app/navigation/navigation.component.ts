import { Component } from '@angular/core';
import { ProductService } from '../product.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-navigation',
  templateUrl: './navigation.component.html',
  styleUrls: ['./navigation.component.css']
})
export class NavigationComponent {
  constructor(private productService: ProductService, private router: Router) {}

  goToInicio(): void {
    this.productService.setProductToEdit(null); 
    this.router.navigate(['/']); 
  }

  goToProductos(): void {
    this.router.navigate(['/productos']);
  }
}
