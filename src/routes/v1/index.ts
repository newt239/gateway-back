import express from 'express';
import authRouter from './auth';
import activityRouter from './activity';
import guestsRouter from './guests';
import exhibitRouter from './exhibit';
import reservationRouter from './reservation';
import adminRouter from './admin';
const router = express.Router();

router.get('/', (_req: express.Request, res: express.Response) => {
    res.json({
        message: 'OK'
    });
});

router.use('/auth/', authRouter);

router.use('/activity/', activityRouter);

router.use('/guests/', guestsRouter);

router.use('/exhibit/', exhibitRouter);

router.use('/reservation/', reservationRouter);

router.use('/admin/', adminRouter);

export default router;