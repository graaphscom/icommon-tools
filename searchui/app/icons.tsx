import {FC, ReactNode} from "react";

export const Icommon: FC<{ node: IcommonNode }> = ({node}) => {
    const Component = node[0] as unknown as FC<Record<string, ReactNode>>

    return <Component {...node[1]}>{node[2]?.map((v, idx) => <Icommon node={v} key={idx}/>)}</Component>
}

export type IcommonNode = [string, Record<string, string | number>, IcommonNode[]?];