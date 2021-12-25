import express from 'express';
import verifyToken from '@/jwt';
var router = express.Router();

router.get('/me', verifyToken, function (req: express.Request, res: express.Response) {
    return res.json({ username: res.locals.username });
});

module.exports = router;