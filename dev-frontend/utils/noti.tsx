import { IconX, IconCheck } from "@tabler/icons-react";
import { v4 } from "uuid";
import { notifications } from "@mantine/notifications";
import { ReactNode } from "react"

type NotifyType = "Success" | "Error"

export function noti(type: NotifyType, title: ReactNode, message: ReactNode) {
    notifications.show({
        id: v4(),
        withCloseButton: true,
        autoClose: 5000,
        title,
        message,
        color: type === "Error" ? 'red' : 'green',
        icon: type === "Error" ? <IconX /> : <IconCheck />,
        loading: false,
        variant: "outline"
    });
}