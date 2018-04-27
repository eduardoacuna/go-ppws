import React, { Component, Fragment } from 'react'
import { Link, Route } from 'react-router-dom'
import { Segment, Header, Container, Button } from 'semantic-ui-react'
import { ToastContainer } from 'react-toastify'
import InteractiveGame from './InteractiveGame'
import ControlPanel from './ControlPanel'

import 'react-toastify/dist/ReactToastify.min.css'
import './App.css'

class App extends Component {
  render() {
    return (
      <Fragment>
        { this.renderToast() }
        { this.renderTitle() }
        { this.renderMenu() }
        { this.renderContent() }
      </Fragment>
    )
  }

  renderToast = () => (
    <ToastContainer
      position='top-right'
      autoClose={8000}
      hideProgressBar
      newestOnTop={false}
      closeOnClick
      pauseOnVisibilityChange
      draggable
      pauseOnHover
    />
  )

  renderTitle() {
    return (
      <Segment vertical textAlign='center' size='big'>
        <Header as='h1' color='red' size='huge'>RoboWars</Header>
      </Segment>
    )
  }

  renderMenu() {
    return (
      <Segment vertical secondary textAlign='center'>
        <Container>
          <Button.Group widths='3'>
            <Button as={Link} color='blue' to='/'>
              Control Panel
            </Button>
            <Button.Or />
            <Button as={Link} color='green' to='/play'>
              Play
            </Button>
            <Button.Or />
            <Button as={Link} color='yellow' to='/watch'>
              Watch
            </Button>
          </Button.Group>
        </Container>
      </Segment>
    )
  }

  renderContent() {
    return (
      <Segment vertical tertiary textAlign='center' className='content'>
        <Route exact path='/' component={Config} />
        <Route path='/play' component={Play} />
        <Route path='/watch' component={Watch} />
      </Segment>
    )
  }
}

const Config = () => (
  <Fragment>
    <Header as='h2'>Control Panel</Header>
    <ControlPanel />
  </Fragment>
)

const Play = () => (
  <Fragment>
    <Header as='h2'>Play</Header>
    <InteractiveGame />
  </Fragment>
)

const Watch = () => (
  <Fragment>
    <Header as='h2'>Watch</Header>
  </Fragment>
)

export default App
