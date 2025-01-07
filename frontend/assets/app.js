// import './app.css'; // Импортируем стили

import { FormController, PopupController } from './javascripts'

class App {
  constructor() {
    this.App = {
      controllers: {
        form: {
          instance: FormController,
          haveElement: true
        },
        popup: {
          instance: PopupController,
          haveElement: false
        },
      }
    }
    window.App = this.App

    this.init()
  }
  init = () => {
    // init controllers
    for (const [name, settings] of Object.entries(this.App.controllers)) {
      const element = document.querySelector(`[data-controller="${name}"]`)

      new settings.instance(settings.haveElement ? element : null)
    }
  }
}

// init App
document.addEventListener('DOMContentLoaded', () => new App())