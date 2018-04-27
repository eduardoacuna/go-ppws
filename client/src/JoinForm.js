import React, { Component } from 'react'
import { Segment, Form, Grid, Message } from 'semantic-ui-react'
import cssColorKeywords from 'css-color-keywords'

class JoinForm extends Component {
  constructor(props) {
    super(props)

    this.colors = Object.keys(cssColorKeywords).map(color => {
      return {
        key: color,
        text: color,
        value: color,
        label: {
          empty: true,
          circular: true,
          style: {
            backgroundColor: color
          }
        }
      }
    })

    this.state = {
      name: '',
      color: ''
    }
  }

  render() {
    const { name, color } = this.state
    const { connecting, error, success } = this.props

    return (
      <Grid stackable columns={3}>
        <Grid.Column width={4}>
        </Grid.Column>
        <Grid.Column width={8}>
          <Segment raised>
            <Form
              onSubmit={this.handleSubmit}
              loading={connecting}
              error={error}
              success={success}
            >
              <Form.Input
                required
                name='name'
                value={name}
                label='Name'
                placeholder='Enter a name'
                onChange={this.handleChange}
                disabled={success}
              />
              <Form.Select
                required
                name='color'
                value={color}
                label='Color'
                options={this.colors}
                placeholder='Select a color'
                onChange={this.handleChange}
                disabled={success}
              />
              <Form.Button
                content='Connect'
                disabled={success}
              />
              <Message
                error
                header='Oops!'
                content='Could not connect to the game server'
              />
              <Message
                success
                header='Please wait'
                content='The game will start shortly'
              />
            </Form>
          </Segment>
        </Grid.Column>
        <Grid.Column width={4}>
        </Grid.Column>
      </Grid>
    )
  }

  handleChange = (e, { name, value }) => this.setState({ [name]: value })

  handleSubmit = () => {
    const { name, color } = this.state
    const { connect } = this.props

    connect(name, color)
  }
}

export default JoinForm