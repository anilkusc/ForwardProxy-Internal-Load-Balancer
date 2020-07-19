import React from 'react';
import RequestCard from '../components/RequestCard';
import Layout from '../components/layout/Layout'


class Home extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      items: [],
    };
  }

  componentDidMount() {
    fetch("http://localhost:8080/api")
      .then(res => res.json())
      .then(
        (result) => {
          this.setState({
            items: result.Logs
          });
        }
      )
  }
  render() {
    return (
      <div>
        <Layout/>
              {//this.state.items.map((data) => <div>{data.Response.Status}</div> )
              this.state.items.map((data) => <RequestCard 
              host={data.Request.Host}
              version ={data.Request.Version}
              method ={data.Request.Method}
              /> )
              
              } 

    </div>
    );
  }
}

export default Home;
