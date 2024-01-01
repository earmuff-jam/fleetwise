import React from 'react';
import moment from 'moment';
import {
  Typography,
  Button,
  Grid,
  Card,
  CardContent,
  makeStyles,
  Box,
  Badge,
  Tooltip,
  CardActions,
  Chip,
  CardMedia,
  Avatar,
} from '@material-ui/core';

const useStyles = makeStyles((theme) => ({
  card: {
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    borderRadius: theme.spacing(1),
  },
  cardMedia: {
    width: theme.spacing(6),
    height: theme.spacing(6),
  },
  contentHeader: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  cardContent: {
    flexGrow: 1,
  },
  headerText: {
    fontSize: '0.925rem',
    fontWeight: 'bold',
    lineHeight: '1.5rem',
    color: theme.palette.primary.main,
  },
  textDetails: {
    fontSize: '0.825rem',
    lineHeight: '1.5rem',
  },
  smallVariant: {},
  rowAlign: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  buttonContainer: {
    borderRadius: theme.spacing(0),
    color: theme.palette.common.black,
    backgroundColor: theme.palette.secondary.light,
    '&:hover': {
      color: theme.palette.common.black,
      boxShadow: '0px 2px 4px rgba(0, 0, 0, 0.2)',
    },
  },
}));

const ViewFilteredEventList = ({ filteredOptions, handleNavigate }) => {
  const classes = useStyles();
  return (
    <>
      {filteredOptions?.map((event) => {
        const eventID = event.id;
        const formattedDate = moment(event.start_date).fromNow();
        const formattedDay = moment(event.start_date).format('dd');
        return (
          <Grid item xs={12} sm={3} key={event.id}>
            <Card className={classes.card}>
              <CardContent className={classes.cardContent}>
                <Box className={classes.contentHeader}>
                  <Box className={classes.rowAlign}>
                    <Avatar
                      alt="Empty image for selected event"
                      className={classes.cardMedia}
                      src={
                        event?.image_url
                          ? event.image_url && `data:image/png;base64,${event.image_url}`
                          : `${process.env.PUBLIC_URL}/blank-canvas.png`
                      }
                    />
                    <Box>
                      <Typography className={classes.headerText}>{event.title}</Typography>
                      <Typography className={classes.textDetails} gutterBottom>
                        {event.cause}
                      </Typography>
                    </Box>
                  </Box>
                  <Tooltip title={`Start Date: ${formattedDate}`}>
                    <Badge badgeContent={formattedDay} color="primary" overlap="rectangular" />
                  </Tooltip>
                </Box>

                <Box className={classes.rowAlign}>
                  <Typography className={classes.textDetails} gutterBottom>
                    {`${event?.skills_required?.map((v) => v).length} active skills`}
                  </Typography>
                  <Chip label={formattedDate} />
                </Box>
                <Typography className={classes.textDetails} gutterBottom>
                  {`${event.street.length ? event.street : 'Unknown location'}`}
                </Typography>
              </CardContent>
              <CardActions>
                <Button onClick={() => handleNavigate(eventID)} className={classes.buttonContainer}>
                  Learn More
                </Button>
              </CardActions>
            </Card>
          </Grid>
        );
      })}
    </>
  );
};

export default ViewFilteredEventList;
