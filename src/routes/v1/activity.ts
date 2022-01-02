import express from 'express';
import verifyToken from '@/jwt';
import { connectDb } from '@/db'
const router = express.Router();

router.post('/:activity_type', verifyToken, function (req: express.Request, res: express.Response) {
    const activity_type: string = req.params.activity_type;
    if (activity_type === "enter" || activity_type === "exit" || activity_type === "pass") {
        const connection = connectDb(res.locals.userid, res.locals.password);
        const guest_id: string = req.body.guest_id;
        const guest_type: string = req.body.guest_type;
        const exhibit_id: string = req.body.exhibit_id;
        const date = new Date();
        const activity_id: string = date.getTime().toString(16) + Math.floor(Math.random() * 10).toString(16);
        const timestamp = date.toLocaleString('ja-JP').slice(0, 19).replace('T', ' ');
        const sql: string = `INSERT INTO gateway.activity (activity_id, guest_id, guest_type, exhibit_id, userid, activity_type, timestamp, available) VALUES ('${activity_id}', '${guest_id}', '${guest_type}', '${exhibit_id}', '${res.locals.userid}', '${activity_type}', '${timestamp}', 1)`;
        connection.query(sql, function (err: any, result: any) {
            if (err) {
                return res.status(400).json(err);
            } else {
                return res.json({
                    status: "success",
                    data: { activity_id: activity_id }
                });
            };
        });
    } else {
        res.json({ status: "error", message: "you posted invaild activity_type" })
    };
});

module.exports = router;