// ** Icon imports
import HomeOutline from 'mdi-material-ui/HomeOutline'
import AccountCogOutline from 'mdi-material-ui/AccountCogOutline'
import AlertCircleOutline from 'mdi-material-ui/AlertCircleOutline'

// ** Type import
import { VerticalNavItemsType } from 'src/@core/layouts/types'

const navigation = (): VerticalNavItemsType => {
  return [
    {
      title: 'Monitor',
      icon: HomeOutline,
      path: '/'
    },
    {
      title: 'Feed Settings',
      icon: AccountCogOutline,
      path: '/feed-settings'
    },
    {
      title: 'Logs',
      icon: AlertCircleOutline,
      path: '/logs',
    }
  ]
}

export default navigation
