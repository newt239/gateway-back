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

const placesRouter = require('./places');
router.use('/places/', placesRouter);

module.exports = router;