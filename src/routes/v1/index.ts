import express from 'express';
const router = express.Router();

router.get('/', function (req: express.Request, res: express.Response) {
    res.json({
        message: 'OK'
    });
});

const authRouter = require('./auth');
router.use('/auth/', authRouter);

const usersRouter = require('./users');
router.use('/users/', usersRouter);

const activityRouter = require('./activity');
router.use('/activity/', activityRouter);

module.exports = router;