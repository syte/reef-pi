import React from 'react'
import { render } from 'react-dom'
import MainPanel from './main_panel.jsx'
import Common from './common.jsx'
import SignIn from './sign_in.jsx'

export default class App extends Common {
  constructor (props) {
    super(props)
    this.state = {
      info: {}
    }
    this.loadInfo = this.loadInfo.bind(this)
  }
  componentDidMount () {
    this.loadInfo()
  }

  loadInfo () {
    this.ajaxGet({
      url: '/api/info',
      success: function (data) {
        this.setState({
          info: data
        })
      }.bind(this)
    })
  }

  render () {
    if (!SignIn.isSignIned()) {
      return (<SignIn />)
    }
    var st = {textAlign: 'center'}
    return (
      <div className='container'>
        {super.render()}
        <div className='container'><h3 style={st}> {this.state.info.name} </h3></div>
        <div className='container'>
          <MainPanel />
        </div>
      </div>
    )
  }
}

render(<App />, document.getElementById('main-panel'))
