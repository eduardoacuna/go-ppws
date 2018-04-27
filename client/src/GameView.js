import React, { Component, Fragment } from 'react'

class GameView extends Component {
  constructor(props) {
    super(props)
    this.state = {
      unit: 0
    }
    this.svgRef = React.createRef()
  }

  componentDidMount() {
    window.addEventListener('resize', this.updateUnit)
    this.updateUnit()
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.updateUnit)
  }

  updateUnit = () => {
    const {
      gameState: {
        grid: {rows, cols}
      }
    } = this.props

    this.svgRef.current.setAttribute('display', 'none')

    const width = this.svgRef.current.parentNode.clientWidth - 50
    const height = this.svgRef.current.parentNode.clientHeight - 70

    this.svgRef.current.setAttribute('display', 'inline')

    let minSize = width < height ? width : height
    let maxDims = cols < rows ? rows : cols

    const newUnit = Math.floor(minSize / maxDims)

    this.setState({
      unit: newUnit
    })
  }

  render() {
    const { unit } = this.state
    const {
      gameState: {
        grid: { rows, cols, cells },
        player,
        enemies
      }
    } = this.props

    const width = cols * unit
    const height = rows * unit
    return (
      <svg
        ref={this.svgRef}
        xmlns='http://www.w3.org/2000/svg'
        version='1.1'
        baseProfile='full'
        width={width}
        height={height}
      >
        <Board
          rows={rows}
          cols={cols}
          cells={cells}
          unit={unit}
        />
        <Player
          x={player.position % cols}
          y={Math.floor(player.position / cols)}
          unit={unit}
          direction={player.direction}
          color={player.color}
          name={player.name}
          score={player.score}
        />
        {enemies.map((enemy, index) => (
          <Enemy
            key={index}
            x={enemy.position % cols}
            y={Math.floor(enemy.position / cols)}
            unit={unit}
            direction={enemy.direction}
            color={enemy.color}
            name={enemy.name}
            score={enemy.score}
          />
        ))}
        <text
          x={unit * (player.position % cols)}
          y={unit * Math.floor(player.position / cols)}
          fill='white'
        >
          {`${player.name} (${player.score})`}
        </text>
        {enemies.map((enemy, index) => (
          <text
            key={index}
            x={unit * (enemy.position % cols)}
            y={unit * Math.floor(enemy.position / cols)}
            fill='white'
          >
            {`${enemy.name} (${enemy.score})`}
          </text>
        ))}
      </svg>
    )
  }
}

const Pointer = ({ x, y, unit, direction }) => {
  const pointerRatio = 1/3 * unit
  const pointerRadius = 1/6 * unit
  let pointerX = x
  let pointerY = y
  switch (direction) {
    case 'north':
      pointerY -= pointerRatio
      break
    case 'south':
      pointerY += pointerRatio
      break
    case 'east':
      pointerX += pointerRatio
      break
    case 'west':
      pointerX -= pointerRatio
      break
    default:
  }

  return (
    <circle
      className='pointer'
      cx={pointerX}
      cy={pointerY}
      r={pointerRadius}
    />
  )
}

const Player = ({ x, y, unit, direction, color, name, score }) => {
  const centerX = unit * x + Math.floor(unit / 2)
  const centerY = unit * y + Math.floor(unit / 2)
  const baseRadius = Math.ceil(unit / 2)
  
  return (
    <Fragment>
      <circle
        className='player'
        cx={centerX}
        cy={centerY}
        r={baseRadius}
        fill={color}
        stroke='black'
      />
      <Pointer
        x={centerX}
        y={centerY}
        unit={unit}
        direction={direction}
      />
    </Fragment>
  )
}

const Enemy = ({ x, y, unit, direction, color, name, score }) => {
  const centerX = unit * x + Math.floor(unit / 2)
  const centerY = unit * y + Math.floor(unit / 2)
  const baseRadius = Math.ceil(unit / 2)
  return (
    <Fragment>
      <circle
        className='enemy'
        cx={centerX}
        cy={centerY}
        r={baseRadius}
        fill={color}
        stroke='black'
      />
      <Pointer
        x={centerX}
        y={centerY}
        unit={unit}
        direction={direction}
      />
    </Fragment>
  )
}

const Board = ({rows, cols, unit, cells}) => {
  return cells.map((mark, index) => (
      <Cell
        key={index}
        size={unit}
        mark={mark}
        x={index % cols}
        y={Math.floor(index / cols)}
      />
    )
  )
}

const Cell = ({ size, mark, x, y }) => {
  const markClass = cssMarkClass(mark)
  return (
    <rect
      x={size * x}
      y={size * y}
      width={size}
      height={size}
      className={'board_cell ' + markClass}
    />
  )
}

const cssMarkClass = (mark) => {
  switch (mark) {
    case '_':
      return 'board_cell_floor'
    case '#':
      return 'board_cell_wall'
    default:
      return ''
  }
}

export default GameView