import React, { Component } from 'react'
import { connect } from 'react-redux'
import { addRecord } from '../actions'

class AddRecord extends Component {
  render() {
    const { dispatch } = this.props
    let input
    return (
      <div>
        <form onSubmit={e => {
          e.preventDefault()
          if (!input.value.trim()) {
            return
          }
          dispatch(addRecord(input.value))
          input.value = ''
        }}>
          <input ref={node => input = node} />
          <button type="submit">
            Add Session
          </button>
        </form>
      </div>
    )
  }
}

export default connect()(AddRecord)
