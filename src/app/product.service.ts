import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
@Injectable({
  providedIn: 'root'
})
export class ProductService {
  private storageKey = 'products';  
  private productToEdit = new BehaviorSubject<any>(null);
  constructor() { }

  // Obtener todos los productos 
  getProducts(): any[] {
    const products = localStorage.getItem(this.storageKey);
    return products ? JSON.parse(products) : [];  
  }

  // Guardar un nuevo producto en LocalStorage
  addProduct(product: any): void {
    const products = this.getProducts();
    const id = new Date().getTime();  
    const newProduct = { id, ...product };
    products.push(newProduct);
    localStorage.setItem(this.storageKey, JSON.stringify(products));
  }

  // Eliminar un producto por su ID
  deleteProduct(id: number): void {
    let products = this.getProducts();
    products = products.filter(product => product.id !== id);  
    localStorage.setItem(this.storageKey, JSON.stringify(products));
  }

  // Actualizar un producto existente
  updateProduct(updatedProduct: any): void {
    const products = this.getProducts();
    const index = products.findIndex(prod => prod.id === updatedProduct.id);
  
    if (index !== -1) {
      products[index] = updatedProduct;
      localStorage.setItem('products', JSON.stringify(products));
    }
  }
  setProductToEdit(product: any): void {
    this.productToEdit.next(product); 
  }
  
  getProductToEdit(): any {
    return this.productToEdit.asObservable();
  }
}
