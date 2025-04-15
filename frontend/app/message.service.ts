import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class MessageService {
  private messageSource = new Subject<string>(); 
  currentMessage = this.messageSource.asObservable(); 

  constructor() {}


  changeMessage(message: string) {
    this.messageSource.next(message);
  }
}
