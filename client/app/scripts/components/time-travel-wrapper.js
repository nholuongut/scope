import React from 'react';
import { connect } from 'react-redux';
import { TimeTravel } from 'nholuongutworks-ui-components';

import { jumpToTime, resumeTime, pauseTimeAtNow } from '../actions/request-actions';

class TimeTravelWrapper extends React.Component {
  handleLiveModeChange = (showingLive) => {
    if (showingLive) {
      this.props.resumeTime();
    } else {
      this.props.pauseTimeAtNow();
    }
  }

  render() {
    return (
      <TimeTravel
        hasLiveMode
        timestamp={this.props.timestamp}
        showingLive={this.props.showingLive}
        onChangeTimestamp={this.props.jumpToTime}
        onChangeLiveMode={this.handleLiveModeChange}
      />
    );
  }
}

function mapStateToProps(state) {
  return {
    showingLive: !state.get('pausedAt'),
    timestamp: state.get('pausedAt'),
  };
}

export default connect(
  mapStateToProps,
  {
    jumpToTime, pauseTimeAtNow, resumeTime
  },
)(TimeTravelWrapper);
