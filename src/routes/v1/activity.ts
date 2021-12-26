import express from 'express';
import verifyToken from '@/jwt';
import { connectDb } from '@/db'
const router = express.Router();

router.get('/enter', verifyToken, function (req: express.Request, res: express.Response) {
});

module.exports = router;