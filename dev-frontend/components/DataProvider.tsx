import { ReactNode, useState, useEffect } from "react"
import { atom } from "signia";
import { DataContext } from "./DataContext"
import useWebSocket from "react-use-websocket";

export let bytesWData = atom<{ x: number, y: number }[]>("bytesData", []);
export let bytesRData = atom<{ x: number, y: number }[]>("bytesData", []);
export let oprWData = atom<{ x: number, y: number }[]>("oprData", []);
export let oprRData = atom<{ x: number, y: number }[]>("oprData", []);

export default function DataProvider({ children }: { children: ReactNode }) {
    let [myData, setMyData] = useState(null);
    let { lastMessage } = useWebSocket("ws://192.168.177.201:30401/api/storage/pvcws");
    useEffect(() => {
        if (lastMessage) {
            let newData = JSON.parse(lastMessage.data) as typeof myData;
            setMyData(newData);
            let dataLength = 20;
            let msgData = newData as any;
            if (msgData.CephPerformance) {
                let read = oprRData.value;
                let write = oprWData.value;
                if (read.length === 0) {
                    read = [{
                        x: 1,
                        y: msgData.CephPerformance.ReadOperationsPerSec
                    }]
                } else {
                    let shouldShift = read.length >= dataLength;
                    read.push({
                        x: read[read.length - 1].x + 1,
                        y: msgData.CephPerformance.ReadOperationsPerSec
                    });
                    shouldShift && read.shift();
                }
                if (write.length === 0) {
                    write = [{
                        x: 1,
                        y: msgData.CephPerformance.WriteOperationPerSec
                    }];
                } else {
                    let shouldShift = write.length >= dataLength;
                    write.push({
                        x: write[write.length - 1].x + 1,
                        y: msgData.CephPerformance.WriteOperationPerSec
                    });
                    shouldShift && write.shift();
                };
                oprRData.set(read);
                oprWData.set(write);

                read = bytesRData.value;
                write = bytesWData.value;
                if (read.length === 0) {
                    read = [{
                        x: 1,
                        y: msgData.CephPerformance.ReadBytesPerSec
                    }];
                } else {
                    let shouldShift = read.length >= dataLength;
                    read.push({
                        x: read[read.length - 1].x + 1,
                        y: msgData.CephPerformance.ReadBytesPerSec
                    });
                    shouldShift && read.shift();
                }
                if (write.length === 0) {
                    write = [{
                        x: 1,
                        y: msgData.CephPerformance.WriteBytesPerSec
                    }];
                } else {
                    let shouldShift = write.length >= dataLength;
                    write.push({
                        x: write[write.length - 1].x + 1,
                        y: msgData.CephPerformance.WriteBytesPerSec
                    });
                    shouldShift && write.shift();
                }
                bytesRData.set(read);
                bytesWData.set(write);
            }
        }
    }, [lastMessage]);
    return (
        <DataContext.Provider value={myData}>
            {children}
        </DataContext.Provider>
    )
}