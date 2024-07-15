import { useQuery } from "react-query"
import Config from "../config"
import { UserResponse } from "../model"
const toJson = (res: Response) => res.json()
const queryUser = () => fetch(Config.userUri, {credentials: 'include'}).then(toJson)
export default () => {
    const { data } = useQuery<UserResponse>(
        ['user'],
        queryUser,
        {staleTime: 5*60*1000}, 
    )
    return data
}