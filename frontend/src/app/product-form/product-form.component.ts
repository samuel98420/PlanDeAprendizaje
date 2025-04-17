import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ProductService } from '../product.service';
import { MessageService } from '../message.service';

@Component({
  selector: 'app-product-form',
  templateUrl: './product-form.component.html',
  styleUrls: ['./product-form.component.css']
})
export class ProductFormComponent implements OnInit {
  @Input() product: any = null;
  @Output() close = new EventEmitter<void>();
  productForm: FormGroup;
  formSubmitted = false;
  isEditMode: boolean = false;
  imagePreview: string | ArrayBuffer | null = null; 
  imageBase64: string = ''; 
  selectedFileName = ''; 


 
  constructor(private productService: ProductService, private fb: FormBuilder, private messageService: MessageService) {
    this.productForm = this.fb.group({
      name: ['', Validators.required],
      price: ['', [Validators.required, Validators.min(0)]],
      description: ['', Validators.required],
      imageUrl: [''] 
    
    });
  }

  ngOnInit(): void {
    if (this.product && this.product.id) {
      this.isEditMode = true;
      this.productForm.patchValue({
        name: this.product.name,
        price: this.product.price,
        description: this.product.description,
        imageUrl: this.product.imageUrl
      });
      this.imagePreview = this.product.imageUrl; 
    } else {
      this.isEditMode = false;
      this.productForm.reset();
    }
  }

  generateRandomId(): string {
    return Math.floor(10000 + Math.random() * 90000).toString();
  }

  onFileSelected(event: any) {
    const file: File = event.target.files[0];
    if (file) {
      this.selectedFileName = file.name;
      const reader = new FileReader();
      reader.onload = (e: any) => {
        this.imagePreview = e.target.result;
        this.resizeImage(e.target.result);
      };
      reader.readAsDataURL(file);
    } else {
      this.selectedFileName = '';
    }
  }
  
  resizeImage(base64Str: string) {
    const img = new Image();
    img.src = base64Str;
    img.onload = () => {
      const canvas = document.createElement("canvas");
      const ctx = canvas.getContext("2d");

      const maxWidth = 150; 
      const maxHeight = 150;
      let width = img.width;
      let height = img.height;

      if (width > height) {
        if (width > maxWidth) {
          height *= maxWidth / width;
          width = maxWidth;
        }
      } else {
        if (height > maxHeight) {
          width *= maxHeight / height;
          height = maxHeight;
        }
      }

      canvas.width = width;
      canvas.height = height;
      ctx?.drawImage(img, 0, 0, width, height);

      this.imageBase64 = canvas.toDataURL("image/jpeg", 0.7);
    };
  }

  removeImage() {
    this.imagePreview = null;
    this.imageBase64 = '';
    this.selectedFileName = '';
    const fileInput = document.getElementById('image') as HTMLInputElement;
    if (fileInput) {
      fileInput.value = '' 
    }
 } 
  
  onSubmit(): void {
    if (this.productForm.valid) {
      let product = this.productForm.value;
      product.imageUrl = this.imageBase64; 

      if (!this.isEditMode) {
        product.id = this.generateRandomId();
      } else {
        product.id = this.product.id;
      }

      if (this.isEditMode && this.product?.id) {
        this.productService.updateProduct(product);
        this.messageService.changeMessage('Product successfully edited');
      } else {
        this.productService.addProduct(product);
        this.messageService.changeMessage('Product successfully added');
      }

      this.formSubmitted = true;
      this.isEditMode = false;
      this.productForm.reset();
      this.imagePreview = null; 
      this.close.emit();
    }
  }
}
