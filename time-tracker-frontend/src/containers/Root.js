import React, { Component } from 'react'
import { Provider } from 'react-redux'
import configureStore from '../configureStore'
import AsyncApp from './AsyncApp'
import AddRecord from './AddRecord'
import VisibleRecordList from './VisibleRecordList'

const store = configureStore()

export default class Root extends Component {
  render() {
    return (
      <Provider store={store}>
        <AsyncApp />
        <AddRecord />
        <VisibleRecordList />
      </Provider>
    )
  }
}


