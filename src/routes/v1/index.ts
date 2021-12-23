import express from 'express';
var router = express.Router();

router.get('/', function (req: express.Request, res: express.Response) {
    res.json({
        message: 'OK'
    });
});

module.exports = router;