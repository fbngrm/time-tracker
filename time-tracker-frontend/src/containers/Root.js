import React, { Component } from 'react'
import { Provider } from 'react-redux'
import configureStore from '../configureStore'
import Records from './Records'
import AddTimer from './AddTimer'

const store = configureStore()

export default class Root extends Component {
  render() {
    return (
      <Provider store={store}>
        <h2>Add Record</h2>
        <AddTimer />
        <h2>Records</h2>
        <Records />
      </Provider>
    )
  }
}


