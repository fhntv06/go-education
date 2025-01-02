import './app.css'; // Импортируем стили

import { FormController, PopupController } from './javascripts'

class App {
  constructor() {
    this.init()
  }
  init = () => {
    // init classes
    new FormController()
    new PopupController()
  }
}

// init App
document.addEventListener('DOMContentLoaded', () => new App())