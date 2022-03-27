import express from 'express';
const router = express.Router();

router.get('/', function (req: express.Request, res: express.Response) {
    res.json({
        message: 'OK'
    });
});

const authRouter = require('./auth');
router.use('/auth/', authRouter);

const activityRouter = require('./activity');
router.use('/activity/', activityRouter);

const guestsRouter = require('./guests');
router.use('/guests/', guestsRouter);

const exhibitRouter = require('./exhibit');
router.use('/exhibit/', exhibitRouter);

const reservationRouter = require('./reservation');
router.use('/reservation/', reservationRouter);

const adminRouter = require('./admin');
router.use('/admin/', adminRouter);

module.exports = router;