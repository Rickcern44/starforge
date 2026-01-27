import type {PageLoad} from "./$types"

export const load: PageLoad = ({params}) => {
    return {
        pageData: {
            leagueId: params.leagueId,
        }
    }
}