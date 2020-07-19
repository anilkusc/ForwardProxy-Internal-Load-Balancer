import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';

class Bar extends React.Component {
    constructor(props) {
        super(props);
    
        this.state = {
          value:0,
          setValue:0,
        };
        this.handleChange = this.handleChange.bind(this);
      }

    handleChange = (event, newValue) => {
        this.setState({
            value: newValue
          });
    };

    render() {
        const classes = makeStyles({
            root: {
                flexGrow: 1,
                //backgroundColor: '#3f50b5'              
            },
        });
        return (
            <div className={classes.root}>
                <Paper className={classes.root}>
                    <Tabs
                        value={this.state.value}
                        onChange={this.handleChange}
                        indicatorColor="primary"
                        textColor="primary"
                        centered
                    >
                        
                        <Tab label="Connections" />
                        <Tab label="Statistic" />
                        <Tab label="Network Card" />
                        
                    </Tabs>
                </Paper>
            </div>
        );
    }
}
export default Bar;