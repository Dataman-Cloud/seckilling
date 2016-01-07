'use strict';

var React = require('react');
var ReactDOM = require('react-dom');
var ReactBootstrap = require('react-bootstrap');
var ChartistGraph = require('react-chartist');
var $ = require('jquery');

var Button = ReactBootstrap.Button;
var Row = ReactBootstrap.Row;
var Col = ReactBootstrap.Col;
var Pagination = ReactBootstrap.Pagination;
var Modal = ReactBootstrap.Modal;

var API_ROOT = 'http://127.0.0.1:8000/rest/v1/';

var statusMapping = {
    'waiting': '未开始',
    'running': '进行中',
    'end': '已结束'
};

var EventsApp = React.createClass({
    getInitialState: function () {
	    var events = [];
        $.ajax({
            url: API_ROOT + 'activities',
            dataType: "json",
            method: "GET",
            success: function (response) {
                for (var idx in response.results) {
                    events.push(response.results[idx]);
		        }
                this.setState({busy: false,
                               events: events});
	        }.bind(this),
            error: function (xhr, status, err) {
                this.displayError(API_ROOT, 'GET', xhr, status, err);
            }.bind(this)
        });
        return {
            count: 0,
            data: events,
            url: API_ROOT + 'activities',
            params: {},
            page: 1,
            busy: false,
            error: {},
            showresult: true,
        };
    },
    displayError: function (url, method, xhr, status, err) {
        console.log(url, status, err);
        this.setState({
            busy: false,
            error: {
                url: url,
                xhr: xhr,
                status: status,
                err: err,
                method: method
            }
        });
    },
    clearError: function () {
        this.setState({error: {}});
    },
    render: function () {
        return (
            <div className="container-fluid">
                <BarChart data={this.state.data}  showresult={this.state.showresult} />
                <Table data={this.state.data}  showresult={this.state.showresult} />
                <Spinner enabled={this.state.busy} />
                <NetworkErrorDialog onClose={this.clearError} data={this.state.error} />
            </div>
        );
    }
});

var NetworkErrorDialog = React.createClass({
    handleClose: function () {
        this.props.onClose();
    },
    render: function () {
        if (Object.keys(this.props.data).length == 0) {
            return <div />;
        }
        var title = "Error: " + this.props.data.xhr.status + " " + this.props.data.err;
        var resp = <span />;
        if (this.props.data.xhr.responseJSON) {
            resp = (
                <div>
                    <p>The response was:</p>
                    <pre>{JSON.stringify(this.props.data.xhr.responseJSON, null, 2)}</pre>
                </div>
            );
        }
        return (
            <Modal
                title={title}
                className="error-dialog"
                onRequestHide={this.handleClose}>
                <div className="modal-body">
                    <p><b>{this.props.data.method}</b> <code>{this.props.data.url}</code></p>
                    {resp}
                </div>
                <div className='modal-footer'>
                    <Button onClick={this.handleClose}>Close</Button>
                </div>
            </Modal>
        );
    }
});

var Spinner = React.createClass({
    render: function () {
        if (this.props.enabled) {
            return (
                <div className='overlay'>
                    <div className='spinner-loader'>Loading …</div>
                </div>
            );
        } else {
            return (<div />);
        }
    }
});


class BarChart extends React.Component {
  render() {
    if (!this.props.showresult) {
      return <div />;
    }
    var labels = [];
    var delivered_counts = [];
    var total_counts = [];
    for (var idx in this.props.data) {
      labels.push(this.props.data[idx].id);
      delivered_counts.push(this.props.data[idx].delivered_prize_count);
      total_counts.push(this.props.data[idx].count);
    }

    var data = {
      labels: labels,
      series: [
        delivered_counts,
        total_counts
      ]
    };

    var options = {
      high: Math.max(...total_counts),
      low: 0,
      showArea: true,
    };

    var type = 'Bar';

    return (
      <div>
        <Col>
          <h3> 活动统计: </h3>
        </Col>
        <div>
          <ChartistGraph data={data} options={options} type={type} />
        </div>
      </div>
    )
  }
}

var Table = React.createClass({
    render: function () {
        if (!this.props.showresult) {
            return <div />;
        }
        return (
            <div>
                <Col>
                    <h3> 活动列表: </h3>
                </Col>
                <Col>
                    <table className="table table-striped table-bordered table-hover">
                        <thead>
                            <tr>
                                <th className="text-center">ID</th>
                                <th className="text-center">品牌</th>
                                <th className="text-center">发出数</th>
                                <th className="text-center">总数</th>
                                <th className="text-center">开始时间</th>
                                <th className="text-center">结束时间</th>
                                <th className="text-center">状态</th>
                            </tr>
                        </thead>
                        <tbody className="text-center">
                            {(this.props.data || []).map(function(result) {
                                return <EventView key={result.id} data={result}/>;
                            })}
                        </tbody>
                    </table>
                </Col>
            </div>
        );
    }
});

var EventView = React.createClass({
    render: function () {
        return (
            <tr>
                <td><p className="form-control-static">{this.props.data.id}</p></td>
                <td><p className="form-control-static">{this.props.data.brand}</p></td>
                <td><p className="form-control-static">{this.props.data.delivered_prize_count}</p></td>
                <td><p className="form-control-static">{this.props.data.count}</p></td>
                <td><p className="form-control-static">{this.props.data.start_at}</p></td>
                <td><p className="form-control-static">{this.props.data.end_at}</p></td>
                <td><p className="form-control-static">{statusMapping[this.props.data.status]}</p></td>
            </tr>
        );
    }
});

ReactDOM.render(
    <EventsApp />,
    document.getElementById('events')
);
