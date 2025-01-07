export class FormController {
  constructor(element) {
    this.form = element

    this.init()
  }
  init = () => {
    this.fields = Array.from(this.form.querySelectorAll('input'))
    this.button = this.form.querySelector('[type="submit"]')

    // execute methods
    this.addListeners()
    this.checkedAllFieldsFilled()
  }
  addListeners = () => {
    this.form.addEventListener('submit', this.submit)
    this.form.addEventListener('input', this.checkedAllFieldsFilled)
  }
  checkedAllFieldsFilled = () => {
    let allFieldsFilled = true

    for (const field of this.fields) {
      if (!field.value.trim()) {
        allFieldsFilled = false
        break
      }
    }

    this.button.disabled = !allFieldsFilled
  }
  clearValues = (fields) => {
    fields.forEach((field) => field.value = '')
  }
  clearErrors = (fields) => {
    fields.forEach((field) => {
      const parentElement = field.parentElement
      const errorElement = field.parentElement.querySelector('.error')

      parentElement.classList.remove('error')

      if (errorElement) {
        errorElement.remove()
      }
    })
  }
  validation = (fields = []) => {
    this.clearErrors(fields)

    const sortFields = fields.filter((field) => {
      switch (field.name) {
        case 'username':
          return field.value.length  === 0
        case 'email':
          return !field.value.includes('@')
        case 'password':
          return field.value.length < 6
        case 'confirm_password':
          return this.fields.filter(({ name }) => (name === 'password'))[0].value !== field.value
      }
    })

    return sortFields.map((field) => {
      const parent = field.parentElement

      switch (field.name) {
        case 'username':
          if (field.value.length  === 0) {
            parent.classList.add('error')
            // return { id: field.id, text: `Некорректные данные. Проверьте поле ${field.name}` }
          }
          break
        case 'email':
          if (!field.value.includes('@')) {
            parent.classList.add('error')
            return { id: field.id, text: `Некорректные данные. Проверьте поле ${field.name}` }
          }
          break
        case 'password':
          if (field.value.length < 6) {
            parent.classList.add('error')
            return { id: field.id, text: `Некорректные данные. Проверьте поле ${field.name}` }
          }
          break
        case 'confirm_password':
          if (this.fields.filter(({ name }) => (name === 'password'))[0].value !== field.value) {
            parent.classList.add('error')
            return { id: field.id, text: `Некорректные данные. Проверьте поле ${field.name}` }
          }
      }
    })
  }
  showNotification = (messages = []) => {
    messages.forEach((message) => {
      if (message) {
        const field = document.getElementById(message.id)

        field.parentElement.classList.add('error')
        field.insertAdjacentHTML('afterend', `<span class="error">${message.text}</span>`)
      }
    })
  }
  executePopup = (detail) => {
    // Вызвать глобальное событие для открытия попапа
    window.dispatchEvent(new CustomEvent('showPopup', { detail }))
  }
  submit = (event) => {
    event.preventDefault()
    const { action, method } = this.form
    const fields = this.fields
    const errors = [] // this.validation(fields)

    if (errors.length) this.showNotification(errors)
    else {
      fetch(action, {
        method: method || 'GET',
        body: new FormData(this.form)
      })
        .then(async (res) => {
          const { id, text, type } = await res.json()

          switch (type) {
            case 'success':
              this.button.disabled = true

              this.clearValues(fields)
              this.clearErrors(fields)

              this.executePopup({ message: text, type } )
              break;
            case 'error':
              this.showNotification([{ id, text }])
              this.executePopup({ message: text, type } )
          }
        })
        .catch((error) => {
          console.error('Error in form: ' + error)
        })
    }
  }
}