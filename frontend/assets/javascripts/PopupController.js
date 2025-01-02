export class PopupController {
  constructor() {
    this.popups = []
    this.timers = []

    window.addEventListener('showPopup', (event) => {
      if (event.detail && event.detail.message) {
        console.log(event.detail)
        this.showPopup(event.detail)
      }
    })
  }
  createPopup(data) {
    let popup = document.createElement('div')
    popup.classList.add('popup')
    popup.classList.add(data.type)
    popup.insertAdjacentHTML('beforeend', `
      <span class="close-button" onclick="this.parentElement.remove()">×</span>
      <p class="popup-message">${data.message}</p>
    `)
    document.body.insertAdjacentElement('beforeend', popup)
    return popup
  }
  showPopup(data) {
    let popup = this.createPopup(data)

    this.timers.push(setTimeout(() => {
      popup.classList.add('show')
      this.popups.push(popup)
    }, 100))

    this.timers.push(setTimeout(() => {
      popup.classList.add('hidden')
      this.removePopup(popup)
    }, 2000))
  }
  removePopup(popup) {
    popup.addEventListener('transitionend', () => {
      // popup.remove()
      this.popups.pop()
      clearTimeout(this.timers.pop())
    })
  }
}
