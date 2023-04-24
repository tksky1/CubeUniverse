import { commonLayout } from "@/utils/commonLayout"
import {
    Button,
    Center,
    Modal,
    NumberInput,
    Stack,
    Switch,
    TextInput,
    Title,
    Tooltip,
    ActionIcon,
    Loader,
    Group,
    Paper,
    Highlight,
    useMantineTheme
} from "@mantine/core";
import { useEffect, useState, useCallback } from "react";
import { checkStorageOpen, createPvc, deletePvc, getPvcInfo, openStorage, updatePvcVolume } from "@/apis";
import { notifications } from "@mantine/notifications"
import { Box } from "@mantine/core";
import { MantineReactTable } from 'mantine-react-table';
import type { MRT_ColumnDef } from 'mantine-react-table';
import { IconX, IconCheck, IconEdit, IconTrash } from "@tabler/icons-react";
import { useDisclosure } from "@mantine/hooks";
import { useForm } from "@mantine/form";
import { v4 } from "uuid";
import fileStorage from "../public/fileStorage.jpg";
import Image from "next/image";

export default function FileStorage() {
    let [storageOpen, setStorageOpen] = useState<"loading" | false | true>("loading");
    let [pvcs, setPvcs] = useState<FilePVC[]>([]);
    let updateStorageOpen = useCallback(() => {
        checkStorageOpen("file")
            .then(e => e.json())
            .then(e => {
                if (e.msg && e.msg[0] === "p") {
                    setStorageOpen(true);
                } else {
                    setStorageOpen(false);
                }
            })
    }, []);
    let updatePvc = useCallback(() => (
        getPvcInfo("file")
            .then(r => r.json())
            .then(r => setPvcs(Object.values(r.data)))
            .catch(r => {
                notifications.show({
                    id: v4(),
                    withCloseButton: true,
                    autoClose: 5000,
                    title: "Error",
                    message: "Something went wrong",
                    color: 'red',
                    icon: <IconX />,
                    loading: false,
                    variant: "outline"
                });
            })), []);
    useEffect(() => {
        updateStorageOpen();
        updatePvc();
    }, [])
    return (
        storageOpen === "loading"
            ? (
                <Center h="100%">
                    <Loader size={"lg"} />
                </Center>
            )
            : storageOpen
                ? (
                    <Box
                        p={"md"}
                        sx={theme => ({
                        })}>
                        <Stack>
                            <FileStorageInfo />
                            <FilePVCTable pvcs={pvcs} afterSubmit={updatePvc} />
                        </Stack>
                    </Box>
                )
                : (
                    <Center
                        h={"100%"}
                    >
                        <Stack>
                            <Title order={3}>
                                Not Open
                            </Title>
                            <Button onClick={() => {
                                openStorage("file")
                                    .then(r => r.json())
                                    .then(() => updateStorageOpen());
                            }}>
                                Open Me
                            </Button>
                        </Stack>
                    </Center>
                )
    )
}
FileStorage.getLayout = commonLayout;

interface FilePVC {
    name: string,
    namespace: string,
    time: string,
    volume: string
}

let hostColumns: MRT_ColumnDef<FilePVC>[] = [
    {
        header: "存储卷声明",
        accessorKey: "name"
    },
    {
        header: " 命名空间",
        accessorKey: "namespace",
    },
    {
        header: "存储卷容量",
        accessorKey: "volume"
    },
    {
        header: "创建时间",
        accessorKey: "time"
    }
]

interface FilePVCTableInterface {
    pvcs: FilePVC[],
    afterSubmit: () => Promise<void>
}

function FilePVCTable({ pvcs, afterSubmit }: FilePVCTableInterface) {
    let form = useForm({
        initialValues: {
            name: "",
            namespace: "",
            volume: 1,
            autoscale: false
        },
        validate: {
            name: val => val.length === 0 ? "Empty Name" : null,
            namespace: val => val.length === 0 ? "Empty Namespace" : null,
            volume: val => (typeof val === "number" && val >= 1) ? null : "Invalid Value"
        }
    });
    let [opened, { open, close }] = useDisclosure(false);
    let [editOpened, { close: editClose, open: editOpen }] = useDisclosure(false);
    let [deleteOpened, { close: deleteClose, open: deleteOpen }] = useDisclosure(false);
    let [currentRow, setCurrentRow] = useState<null | { name: string, namespace: string }>(null);
    useEffect(() => {
        form.reset();
    }, [opened]);
    return (
        <>
            <MantineReactTable
                columns={hostColumns}
                data={pvcs}
                enableColumnOrdering
                enableGlobalFilter
                enableBottomToolbar
                renderTopToolbarCustomActions={() => {
                    return (
                        <Button
                            variant="gradient"
                            gradient={{ from: 'teal', to: 'lime', deg: 105 }}
                            onClick={open}>
                            创建PVC
                        </Button>)
                }}
                onEditingRowSave={() => { }}
                enableEditing
                renderRowActions={({ row, table }) => (
                    <Box sx={{ display: 'flex', gap: '16px' }}>
                        <Tooltip withArrow position="left" label="编辑">
                            <ActionIcon
                                onClick={() => {
                                    setCurrentRow({
                                        name: row.original.name,
                                        namespace: row.original.namespace
                                    });
                                    editOpen();
                                }}
                            >
                                <IconEdit />
                            </ActionIcon>
                        </Tooltip>
                        <Tooltip withArrow position="right" label="删除">
                            <ActionIcon
                                color="red"
                                onClick={() => {
                                    setCurrentRow({
                                        name: row.original.name,
                                        namespace: row.original.namespace
                                    });
                                    deleteOpen();
                                }}>
                                <IconTrash />
                            </ActionIcon>
                        </Tooltip>
                    </Box>
                )}
            />
            <Modal
                opened={opened}
                onClose={close}
                centered
                title={<Title order={4}>新建PVC</Title>}
                radius={"lg"}
                withCloseButton={false}>
                <form onSubmit={form.onSubmit(values => {
                    createPvc({
                        name: values.name,
                        "X-type": "file",
                        autoscale: values.autoscale ? "true" : "false",
                        namespace: values.namespace,
                        volume: `${values.volume}`
                    }).then(e => e.json())
                        .then((e) => {
                            return e.code === 200
                                ? Promise.resolve(e)
                                : Promise.reject(e);
                        })
                        .then(() => afterSubmit().then(() => {
                            close();
                            notifications.show({
                                id: v4(),
                                withCloseButton: true,
                                autoClose: 5000,
                                title: "Success",
                                message: "Create new PVC",
                                color: "green",
                                icon: <IconCheck />,
                                loading: false,
                            });
                        })).catch(e => {
                            close();
                            notifications.show({
                                id: v4(),
                                withCloseButton: true,
                                autoClose: 5000,
                                title: "Error",
                                message: e.msg,
                                color: "red",
                                icon: <IconX />,
                                loading: false,
                            });
                        });
                })}>
                    <Stack>
                        <TextInput
                            withAsterisk
                            label="存储卷声明"
                            {...form.getInputProps("name")} />
                        <TextInput
                            withAsterisk
                            label="命名空间"
                            {...form.getInputProps("namespace")} />
                        <NumberInput
                            min={1}
                            label="存储卷容量"
                            {...form.getInputProps("volume")} />
                        <Switch
                            label="自动扩容"
                            // checked={autoscale}
                            {...form.getInputProps("autoscale")} />
                        <Button
                            variant="gradient"
                            gradient={{ from: 'teal', to: 'blue', deg: 105 }}
                            type="submit">
                            确定
                        </Button>
                    </Stack>
                </form>
            </Modal>
            <EditModal
                opened={editOpened}
                close={editClose}
                name={currentRow ? currentRow.name : ""}
                namespace={currentRow ? currentRow.namespace : ""}
                afterSubmit={afterSubmit}
            />
            <DeleteModal
                afterSubmit={afterSubmit}
                opened={deleteOpened}
                close={deleteClose}
                name={currentRow ? currentRow.name : ""}
                namespace={currentRow ? currentRow.namespace : ""}
            />
        </>
    )
}

interface EditModalInterface {
    name: string,
    namespace: string,
    afterSubmit: () => Promise<void>,
    opened: boolean,
    close: VoidFunction
}

function EditModal({ name, namespace, afterSubmit, opened, close }: EditModalInterface) {
    let form = useForm({
        initialValues: {
            volume: 1,
        },
        validate: {
            volume: val => (typeof val === "number" && val >= 1) ? null : "Invalid Value"
        }
    });
    useEffect(() => {
        form.reset();
    }, [opened]);
    return (
        <Modal
            opened={opened}
            onClose={close}
            centered
            title={<Title order={4}>编辑</Title>}
            size={"xs"}
            radius={"lg"}
            withCloseButton={false}>
            <form onSubmit={form.onSubmit(values => {
                updatePvcVolume({
                    name: name,
                    "X-type": "file",
                    namespace: namespace,
                    volume: `${values.volume}`
                }).then(e => e.json())
                    .then((e) => {
                        return e.code === 200
                            ? Promise.resolve(e)
                            : Promise.reject(e);
                    })
                    .then(() => afterSubmit().then(() => {
                        close();
                        notifications.show({
                            id: v4(),
                            withCloseButton: true,
                            autoClose: 5000,
                            title: "Success",
                            message: "Update PVC colume",
                            color: "green",
                            icon: <IconCheck />,
                            loading: false,
                        });
                    })).catch(e => {
                        close();
                        notifications.show({
                            id: v4(),
                            withCloseButton: true,
                            autoClose: 5000,
                            title: "Error",
                            message: e.msg,
                            color: "red",
                            icon: <IconX />,
                            loading: false,
                        });
                    });
            })}>
                <Stack>
                    <NumberInput
                        size="md"
                        min={1}
                        label="Volume"
                        {...form.getInputProps("volume")} />
                    <Button
                        variant="gradient"
                        gradient={{ from: 'teal', to: 'blue', deg: 105 }}
                        type="submit">
                        确定
                    </Button>
                </Stack>
            </form>
        </Modal>
    )
}

interface DeleteModalInterface {
    name: string,
    namespace: string,
    afterSubmit: () => Promise<void>,
    opened: boolean,
    close: VoidFunction
}

function DeleteModal({ name, namespace, afterSubmit, opened, close }: DeleteModalInterface) {
    return (
        <Modal
            opened={opened}
            onClose={close}
            centered
            size={"auto"}
            radius={"md"}
            withCloseButton={false}>
            <Stack>
                <Group position="apart" w="100%">
                    <Title order={5}>确定要删除此项PVC吗？</Title>
                    <Group spacing={"xs"}>
                        <Button size="xs" color="gray" onClick={close}>取消</Button>
                        <Button size="xs" color="red" onClick={() => {
                            deletePvc({ name, namespace, "X-type": "file" })
                                .then(e => e.json())
                                .then(e => e.code === 200 ? Promise.resolve(e) : Promise.reject(e))
                                .then(e => {
                                    notifications.show({
                                        id: v4(),
                                        withCloseButton: true,
                                        autoClose: 5000,
                                        title: "Success",
                                        message: "PVC deleted",
                                        color: "green",
                                        icon: <IconCheck />,
                                        loading: false,
                                    });
                                })
                                .catch(e => {
                                    notifications.show({
                                        id: v4(),
                                        withCloseButton: true,
                                        autoClose: 5000,
                                        title: "Error",
                                        message: e.msg,
                                        color: 'red',
                                        icon: <IconX />,
                                        loading: false,
                                        variant: "outline"
                                    });
                                })
                                .finally(() => {
                                    close();
                                    afterSubmit();
                                })
                        }}>删除</Button>
                    </Group>
                </Group>
            </Stack>
        </Modal>
    )
}

function FileStorageInfo() {
    let theme = useMantineTheme();
    return (
        <Center>
            <Paper
                p={"xl"}
                sx={theme => ({
                    borderRadius: 1,
                    boxShadow: "0px 2px 5px 1px rgba(74,122,222,0.75)",
                    border: "1px solid #c0e7ff"
                })}>
                <Stack>
                    <Box mr={100} ml={30} mt={20}>
                        <Image
                            height={200}
                            src={fileStorage}
                            style={{
                                float: "left",
                            }}
                            alt="FileStorage" />
                        <Title
                            align="left"
                            order={3}
                            mb={"md"}
                        >文件存储</Title>
                        <Highlight
                            size={"md"}
                            highlight={['文件存储', 'Kubernetes', "持久化卷"]}
                            highlightStyles={(theme) => ({
                                backgroundImage: theme.fn.linearGradient(45, theme.colors.cyan[5], theme.colors.indigo[5]),
                                fontWeight: 700,
                                WebkitBackgroundClip: 'text',
                                WebkitTextFillColor: 'transparent',
                            })}>
                            文件存储是一种基于文件元数据和二进制数据组织形成的由文件系统维护的存储形式。由于有统一的文件系统维护，挂载后容器不需要考虑如何管理磁盘空间，并且支持多个应用同时使用一片空间。在 Kubernetes 中以持久化卷的形式提供，常用于需要多个容器应用共享一个持久化卷的情况。在控制面板创建PVC后，为集群应用pod绑定该PVC即可挂载。
                        </Highlight>
                    </Box>
                </Stack>
            </Paper>
        </Center>
    )
}