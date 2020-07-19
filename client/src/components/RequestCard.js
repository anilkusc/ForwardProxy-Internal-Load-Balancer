import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import Divider from '@material-ui/core/Divider';


class RequestCard extends React.Component {

  render() {
    const classes = makeStyles({
        root: {
          minWidth: 275,
        },
        bullet: {
          display: 'inline-block',
          margin: '0 2px',
          transform: 'scale(0.8)',
        },
        title: {
          fontSize: 14,
        },
        pos: {
          marginBottom: 12,
        },
      });
    return (
      <div>
        <Divider />
       <Card className={classes.root}>
      <CardContent>
        <Typography className={classes.title} color="textSecondary" gutterBottom>
          Request
        </Typography>
        <Typography variant="h5" component="h2">
          {this.props.host}
        </Typography>
        <Typography className={classes.pos} color="textSecondary">
        {this.props.version} - {this.props.method}
        </Typography>
      </CardContent>
      <CardActions>
        <Button size="small">Body</Button><Button size="small">Headers</Button>
      </CardActions>
    </Card>
    <Divider />
    </div>
    );
  }
}

export default RequestCard;
