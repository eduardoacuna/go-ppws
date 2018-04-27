import React, { Component } from 'react'
import { Segment, Form, Grid, Message } from 'semantic-ui-react'

class ControlPanel extends Component {
  constructor(props) {
    super(props)

    this.state = {
      sendingForm: false,
      error: false,
      success: false,

      rows: 50,
      cols: 50,
      wallRoots: 50,
      wallBuildingProb: 0.9,
      povRadius: 12
    }
  }

  render() {
    const {
      sendingForm, error, success, rows, cols, wallRoots, wallBuildingProb, povRadius
    } = this.state
    return (
      <Grid stackable columns={3}>
        <Grid.Column width={4}>
        </Grid.Column>
        <Grid.Column width={8}>
          <Segment raised>
            <Form
              onSubmit={this.handleSubmit}
              loading={sendingForm}
              error={error}
              success={success}
            >
              <Form.Group widths='equal'>
                <Form.Input
                  required
                  type='number'
                  min={10}
                  max={100}
                  name='rows'
                  value={rows}
                  label='Rows'
                  onChange={this.handleChange}
                  disabled={success}
                />
                <Form.Input
                  required
                  type='number'
                  min={10}
                  max={100}
                  name='cols'
                  value={cols}
                  label='Columns'
                  onChange={this.handleChange}
                  disabled={success}
                />
              </Form.Group>
              <Form.Group widths='equal'>
                <Form.Input
                  required
                  type='number'
                  min={1}
                  max={rows * cols}
                  name='wallRoots'
                  value={wallRoots}
                  label='Wall Roots'
                  onChange={this.handleChange}
                  disabled={success}
                />
                <Form.Input
                  required
                  type='number'
                  name='wallBuildingProb'
                  min={0.0}
                  max={0.9}
                  step={0.01}
                  value={wallBuildingProb}
                  label='Wall Building Probability'
                  onChange={this.handleChange}
                  disabled={success}
                />
              </Form.Group>
              <Form.Input
                required
                type='number'
                min={1}
                max={Math.ceil(Math.sqrt(rows*cols))}
                width={8}
                name='povRadius'
                value={povRadius}
                label='POV Radius'
                onChange={this.handleChange}
                disabled={success}
              />
              <Form.Button
                content={success ? 'Stop Game' : 'Start Game' }
              />
              <Message
                error
                header='Oops!'
                content='Something went wrong'
              />
              <Message
                success
                header='Success!'
                content='The game has started'
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
    let {
      rows, cols, wallRoots, wallBuildingProb, povRadius, success
    } = this.state

    if (success) {
      fetch(`http://localhost:3001/stop-game?secret=cierratesesamo`)
        .then(res => {
          if (res.ok) {
            this.setState({ success: false, error: false })
          }
        })
      return
    }

    rows = encodeURIComponent(rows)
    cols = encodeURIComponent(cols)
    wallRoots = encodeURIComponent(wallRoots)
    wallBuildingProb = encodeURIComponent(wallBuildingProb)
    povRadius = encodeURIComponent(povRadius)

    fetch(`http://localhost:3001/start-game?secret=abretesesamo&rows=${rows}&cols=${cols}&wall-roots=${wallRoots}&wall-building-prob=${wallBuildingProb}&pov-radius=${povRadius}`)
      .then(res => {
        if (!res.ok) {
          this.setState({ error: true, success: false })
        } else {
          this.setState({ success: true, error: false })
        }
      })
  }
}

export default ControlPanel