import { FormController } from './javascripts'

class App {
  constructor() {
    this.init()
  }
  init = () => {
    // init classes
    new FormController()
  }
}

// init App
document.addEventListener('DOMContentLoaded', () => new App())