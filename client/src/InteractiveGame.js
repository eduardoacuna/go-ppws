import React, { Component } from 'react'
import { toast } from 'react-toastify'
import JoinForm from './JoinForm'
import GameView from './GameView'

class InteractiveGame extends Component {
  constructor(props) {
    super(props)

    this.state = {
      wsocket: null,
      gameState: null,
      success: false,
      error: false,
      connecting: false
    }
  }

  componentWillUnmount() {
    const { wsocket } = this.state

    window.removeEventListener('keypress', this.handleKeyPress)

    if (!wsocket) {
      return
    }
    wsocket.close()
  }

  connect = (name, color) => {
    this.setState({connecting: true}, () => {
      name = encodeURIComponent(name)
      color = encodeURIComponent(color)

      const ws = new WebSocket(`ws://localhost:3001/play?name=${name}&color=${color}`)

      ws.onopen = this.handleOpen
      ws.onclose = this.handleClose
      ws.onmessage = this.handleMessage
      ws.onerror = this.handleError

      this.setState({ wsocket: ws, connecting: false }, () => {
        window.addEventListener('keypress', this.handleKeyPress)
      })
    })
  }

  render() {
    const { gameState, connecting, error, success } = this.state

    if (!gameState) {
      return (
        <JoinForm connect={this.connect} connecting={connecting} error={error} success={success} />
      )
    } else {
      return (
        <GameView gameState={gameState} />
      )
    }
  }

  handleOpen = (evt) => {
    this.setState({ success: true, error: false }, () => {
      toast.info('Connection opened', {
        closeOnClick: false,
        draggable: false
      })
    })
  }

  handleClose = (evt) => {
    this.setState({ wsocket: null, gameState: null, success: false, error: false }, () => {
      toast.warn('Connection closed')
    })
  }

  handleMessage = (evt) => {
    this.setState({ gameState: JSON.parse(evt.data) })
  }

  handleError = (evt) => {
    this.setState({ wsocket: null, gameState: null, error: true, success: false }, () => {
      toast.error('Error')
    })
  }

  sendAction = (command) => {
    const { wsocket } = this.state

    wsocket.send(JSON.stringify({command}))
  }

  handleKeyPress = (evt) => {
    console.log(evt)
    const [keyLeft, keyUp, keyRight, keyDown] = [37, 38, 39, 40]

    switch (evt.keyCode) {
      case keyLeft:
        this.sendAction('turn-left')
        break
      case keyUp:
        this.sendAction('move-forward')
        break
      case keyRight:
        this.sendAction('turn-right')
        break
      case keyDown:
        this.sendAction('move-backward')
        break
      default:
        if (evt.charCode === 32) {
          this.sendAction('attack')
        }
    }
  }
}
export default InteractiveGame