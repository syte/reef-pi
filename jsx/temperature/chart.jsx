import React from 'react'
import {ComposedChart, Line, Tooltip, YAxis, XAxis, Bar, ReferenceLine} from 'recharts'
import Common from '../common.jsx'
import $ from 'jquery'

export default class Chart extends Common {
  constructor (props) {
    super(props)
    this.state = {
      usage: []
    }
    this.fetch = this.fetch.bind(this)
  }

  componentDidMount () {
    var timer = window.setInterval(this.fetch, 10 * 1000)
    this.setState({timer: timer})
    this.fetch()
  }

  componentWillUnmount () {
    window.clearInterval(this.state.timer)
  }

  fetch () {
    this.ajaxGet({
      url: '/api/tc/usage',
      success: function (data) {
        var processed = []
        $.each(data, function (i, v) {
          v.cooler *= -1
          processed.push(v)
        })
        this.setState({
          usage: processed,
          showAlert: false
        })
      }.bind(this)
    })
  }

  render () {
    if (this.state.usage.length <= 0) {
      return (<div />)
    }
    return (
      <div className='container'>
        {super.render()}
        <span className='h6'>Heater/Cooler</span>
        <ComposedChart width={500} height={250} data={this.state.usage}>
          <YAxis yAxisId='left' orientation='left' domain={[76, 82]} />
          <YAxis yAxisId='right' orientation='right' />
          <ReferenceLine yAxisId='right' y={0} />
          <XAxis dataKey='time' />
          <Tooltip />
          <Bar dataKey='heater' fill='#ffbb33' isAnimationActive={false} yAxisId='right' stackId='t' />
          <Bar dataKey='cooler' fill='#33b5e5' isAnimationActive={false} yAxisId='right' stackId='t' />
          <Line type='monotone' dataKey='temperature' stroke='#ce93d8' isAnimationActive={false} yAxisId='left' dot={false} />
        </ComposedChart>
      </div>
    )
  }
}
