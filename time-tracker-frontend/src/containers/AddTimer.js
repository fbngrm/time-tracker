import { connect } from 'react-redux'
import { addTimer, startTimer, stopTimer, saveTimerIfNeeded } from '../actions'
import AddTimer from '../components/AddTimer'

function mapStateToProps(state){
  return {
    timer: state.timer
  }
}

function mapDispatchToProps(dispatch){
  return {
    addTimer: function(state){
      return dispatch(addTimer(state))
    },
    startTimer: function(state){
      return dispatch(startTimer(state))
    },
    stopTimer: function(state){
      return dispatch(stopTimer(state))
    },
    saveTimer: function(state){
      return dispatch(saveTimerIfNeeded(state))
    }
  }
}

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(AddTimer)
