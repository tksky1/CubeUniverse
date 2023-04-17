import type { NextApiRequest, NextApiResponse } from "next";

let getData = {
    code: 200,
    data: {
        "block0": {
            "name": "blocktestname-1",
            "namespace": "blocktestnamespace-1",
            "volume": "1",
            "time": "2006-01-02 15:04:00"
        },
        "block1": {
            "name": "blocktestname-2",
            "namespace": "blocktestnamespace-2",
            "volume": "2",
            "time": "2006-01-02 15:04:01"
        },
        "block2": {
            "name": "blocktestname-3",
            "namespace": "blocktestnamespace-3",
            "volume": "3",
            "time": "2006-01-02 15:04:02"
        },
        "block3": {
            "name": "blocktestname-4",
            "namespace": "blocktestnamespace-4",
            "volume": "4",
            "time": "2006-01-02 15:04:03"
        },
        "block4": {
            "name": "blocktestname-5",
            "namespace": "blocktestnamespace-5",
            "volume": "5",
            "time": "2006-01-02 15:04:04"
        },
        "block5": {
            "name": "blocktestname-6",
            "namespace": "blocktestnamespace-6",
            "volume": "6",
            "time": "2006-01-02 15:04:05"
        },
        "block6": {
            "name": "blocktestname-7",
            "namespace": "blocktestnamespace-7",
            "volume": "7",
            "time": "2006-01-02 15:04:06"
        },
        "block7": {
            "name": "blocktestname-8",
            "namespace": "blocktestnamespace-8",
            "volume": "8",
            "time": "2006-01-02 15:04:07"
        },
        "block8": {
            "name": "blocktestname-9",
            "namespace": "blocktestnamespace-9",
            "volume": "9",
            "time": "2006-01-02 15:04:08"
        },
        "block9": {
            "name": "blocktestname-:",
            "namespace": "blocktestnamespace-:",
            "volume": ":",
            "time": "2006-01-02 15:04:09"
        }
    },
    msg: "all info"
}

export default function handler<T extends NextApiRequest>(req: T, res: NextApiResponse) {
    if (req.method === "GET") {
        res.status(200).json(getData);
    } else if (req.method === "DELETE") {
        res.status(200).json({
            "code": 200,
            "data": null,
            "msg": "delete done"
        });
    } else if (req.method === "POST") {
        res.status(200).json({ "code": 200, "data": null, "msg": "create done" })
    }
}