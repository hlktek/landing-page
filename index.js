require('dotenv').config({ path: `${__dirname}/../.env` });
const express = require('express');
const path = require('path');
const bodyParser = require('body-parser');
const moment = require('moment-timezone');
const passport = require('passport');
const config = require('./config');

// const indexRouter = require('./routes');

const port= config.PORT || 3000;
const CONTEXT_PATH = process.env.CONTEXT_PATH || '/';

const app = express();
require('./lib/passport');


moment.tz.setDefault('Asia/Ho_Chi_Minh');

// view engine setup
app.set('views', path.join(__dirname, 'views'));
app.set('view engine', 'ejs');
app.disable('view cache');
app.use(bodyParser.json()); // to support JSON-encoded bodies
app.use(bodyParser.urlencoded({ extended: true }));
app.use(CONTEXT_PATH, express.static(path.join(__dirname, '../public')));
app.get('/', function(req, res) {
  res.render('auth');
});

app.use(passport.initialize());
app.use(passport.session());
require('./router/authRoutes')(app);


app.listen(port, () => {
  console.log(`Server running on port ${port}`);
});
