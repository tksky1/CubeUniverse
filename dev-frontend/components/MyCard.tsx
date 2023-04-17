import { 
    Text, 
    Card, 
    Center
} from "@mantine/core"
import Link from "next/link";
import { ReactElement } from "react";

interface MyCard {
    title: string,
    children: ReactElement,
    href?: string
}

export default function MyCard({title, children, href}: MyCard) {
    return (
        <Card 
            shadow={"xs"}
            p={"xl"}
            h={"100%"}
            // withBorder
            sx={theme => ({
                backgroundColor: "#f7f8fb"
                // backgroundColor: theme.fn.darken(theme.colors.indigo[9], 0.85),
                // borderColor: theme.fn.darken(theme.colors.indigo[9], 0.85)
            })}
        >
            <Card.Section 
                mb={"xs"} 
                sx={_theme => ({
                    borderBottom: "solid 1px rgba(0, 0, 0, 0.3)"
                })}>
                {href
                    ? (<Link href={href}>
                        <Text 
                            size={"xs"}
                            weight={"bold"}
                            underline
                            color="#0e201a"
                        >
                            {title}
                        </Text>
                    </Link>)
                    : (<Text 
                        size={"xs"}
                        weight={"bold"}
                        color="#0e201a"
                    >
                        {title}
                    </Text>)}
            </Card.Section>
            <Center h={"100%"}>
                {children}
            </Center>
        </Card>
    )
}